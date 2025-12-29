#!/usr/bin/env python3
"""
Model Template Config Translator - Python Implementation

Translates model configuration JSON to inference engine CLI arguments
using the declarative mapping files in mappings/.

Usage:
    python3 examples/translate_config.py templates/qwen3-sglang.json
"""

import json
import sys
from pathlib import Path


def load_mappings(engine, mappings_dir):
    """Load engine-specific and common mapping files."""
    engine_path = mappings_dir / f"{engine}.json"
    common_path = mappings_dir / "common.json"

    with open(engine_path) as f:
        engine_mappings = json.load(f)
    with open(common_path) as f:
        common_mappings = json.load(f)

    return engine_mappings, common_mappings


def get_nested_value(data, path):
    """Extract value from nested dict using dot notation."""
    parts = path.split('.')
    current = data

    for part in parts:
        if not isinstance(current, dict) or part not in current:
            return None
        current = current[part]

    return current


def format_value(value, mapping_entry):
    """Format value according to its type specification."""
    if value is None:
        return None, True  # (value, skip)

    entry_type = mapping_entry['type']

    if entry_type == 'string':
        return str(value), False

    elif entry_type == 'int':
        return str(int(value)), False

    elif entry_type == 'float':
        format_str = mapping_entry.get('format', '{0}')
        if format_str.startswith('%.'):
            # Convert C-style format to Python format
            precision = int(format_str[2:-1])
            return f"{float(value):.{precision}f}", False
        return str(float(value)), False

    elif entry_type == 'bool':
        return str(value).lower(), False

    elif entry_type == 'bool_flag':
        # Only include flag if value is true
        if value:
            return None, False
        return None, True

    elif entry_type == 'json_object':
        # Check for empty dict/object
        if isinstance(value, dict) and len(value) == 0:
            return None, True
        return json.dumps(value, separators=(',', ':')), False

    elif entry_type == 'array_of_strings':
        if isinstance(value, list):
            return ' '.join(str(v) for v in value), False
        return str(value), False

    elif entry_type == 'object':
        return json.dumps(value, separators=(',', ':')), False

    else:
        return str(value), False


def get_cli_flag(mapping_entry, engine):
    """Get CLI flag for the specified engine."""
    cli_flag = mapping_entry['cli_flag']

    # Handle engine-specific flags
    if isinstance(cli_flag, dict):
        return cli_flag.get(engine, '')

    return cli_flag


def build_custom_command(config):
    """Build command for custom engine."""
    custom = config.get('custom', {})

    if not custom or 'base_command' not in custom:
        raise ValueError("Custom engine requires 'base_command'")

    # Split base command
    parts = custom['base_command'].split()

    # Add model
    if config.get('model'):
        model_flag = custom.get('model_flag', '--model')
        parts.extend([model_flag, config['model']])

    # Add args
    for arg in custom.get('args', []):
        parts.extend(arg.split())

    # Add kv_args
    for key, value in custom.get('kv_args', {}).items():
        parts.extend([key, str(value)])

    return parts


def build_command(config, engine_mappings, common_mappings, engine):
    """Build CLI command from config and mappings."""
    if engine == 'custom':
        return build_custom_command(config)

    parts = []

    # Add base command
    if 'base_command' in engine_mappings:
        parts.extend(engine_mappings['base_command'])

    # Add model
    if config.get('model'):
        parts.extend(['--model', config['model']])

    # Process engine-specific mappings
    for entry_name, entry in engine_mappings.get('mappings', {}).items():
        if entry['json_path'] == 'model':
            continue  # Already handled

        value = get_nested_value(config, entry['json_path'])
        if value is None:
            continue

        cli_flag = get_cli_flag(entry, engine)
        if not cli_flag:
            continue

        formatted, skip = format_value(value, entry)
        if skip:
            continue

        if entry['type'] == 'bool_flag':
            parts.append(cli_flag)
        else:
            parts.extend([cli_flag, formatted])

    # Process common mappings
    for entry_name, entry in common_mappings.get('mappings', {}).items():
        if entry['json_path'] == 'model':
            continue  # Already handled

        value = get_nested_value(config, entry['json_path'])
        if value is None:
            continue

        cli_flag = get_cli_flag(entry, engine)
        if not cli_flag:
            continue

        formatted, skip = format_value(value, entry)
        if skip:
            continue

        if entry['type'] == 'bool_flag':
            parts.append(cli_flag)
        else:
            parts.extend([cli_flag, formatted])

    return parts


def main():
    if len(sys.argv) < 2:
        print("Usage: python3 examples/translate_config.py <config-file.json>")
        print("\nExample:")
        print("  python3 examples/translate_config.py templates/qwen3-sglang.json")
        sys.exit(1)

    config_path = Path(sys.argv[1])

    # Load config
    with open(config_path) as f:
        config = json.load(f)

    engine = config.get('engine')
    if not engine:
        print("Error: Config must specify an 'engine' field")
        sys.exit(1)

    # Get project root and mappings directory
    script_dir = Path(__file__).parent
    project_root = script_dir.parent
    mappings_dir = project_root / 'mappings'

    print(f"Loading mappings for engine: {engine}")

    # Load mappings
    engine_mappings, common_mappings = load_mappings(engine, mappings_dir)

    # Build command
    command = build_command(config, engine_mappings, common_mappings, engine)

    # Display results
    print("\nGenerated Command:")
    print(' '.join(command))

    print("\nFormatted Command:")
    if command:
        print(command[0], end='')
        for arg in command[1:]:
            if arg.startswith('-'):
                print(f' \\\n  {arg}', end='')
            else:
                print(f' {arg}', end='')
        print()


if __name__ == '__main__':
    main()
