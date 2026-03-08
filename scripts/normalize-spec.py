#!/usr/bin/env python3
"""
Normalize UniFi OpenAPI spec for use with tfplugingen-openapi.
- Renames schema components with spaces to snake_case
- Updates all $ref references to match
- Writes normalized spec to output path
"""

import json
import re
import sys

def to_snake_case(name):
    """Convert a schema name to a valid identifier (snake_case)."""
    # Replace spaces and special chars with underscores
    name = re.sub(r'[^a-zA-Z0-9]', '_', name)
    # Collapse multiple underscores
    name = re.sub(r'_+', '_', name)
    return name.strip('_')

def normalize_spec(input_path, output_path):
    with open(input_path) as f:
        spec = json.load(f)

    schemas = spec.get('components', {}).get('schemas', {})
    
    # Build rename map: old name -> new name
    rename_map = {}
    for name in schemas:
        new_name = to_snake_case(name)
        if new_name != name:
            rename_map[name] = new_name
    
    print(f"Renaming {len(rename_map)} schemas with spaces/special chars", file=sys.stderr)
    
    # Rename schema keys
    new_schemas = {}
    for name, schema in schemas.items():
        new_name = rename_map.get(name, name)
        new_schemas[new_name] = schema
    spec['components']['schemas'] = new_schemas
    
    # Update all $ref occurrences in the entire spec
    spec_str = json.dumps(spec)
    for old, new in rename_map.items():
        old_ref = f'"$ref": "#/components/schemas/{old}"'
        new_ref = f'"$ref": "#/components/schemas/{new}"'
        spec_str = spec_str.replace(old_ref, new_ref)
    
    spec = json.loads(spec_str)
    
    with open(output_path, 'w') as f:
        json.dump(spec, f, indent=2)
    
    print(f"Wrote normalized spec to {output_path}", file=sys.stderr)

if __name__ == '__main__':
    if len(sys.argv) != 3:
        print(f"Usage: {sys.argv[0]} <input.json> <output.json>")
        sys.exit(1)
    normalize_spec(sys.argv[1], sys.argv[2])
