# Copyright (c) 2024 Inetum Poland.
"""
Parser for GITHUB_OUTPUT file format.
"""

from pathlib import Path


def parse_github_output(output_file: Path) -> dict[str, str]:
    """
    Parse a GITHUB_OUTPUT file into a dictionary.
    
    Handles both single-line format:
        key=value
    
    And multiline format:
        key<<EOF
        line1
        line2
        EOF
    
    Args:
        output_file: Path to the GITHUB_OUTPUT file
    
    Returns:
        Dictionary of output key-value pairs
    """
    if not output_file.exists():
        return {}
    
    content = output_file.read_text()
    outputs = {}
    
    lines = content.split('\n')
    i = 0
    
    while i < len(lines):
        line = lines[i]
        
        if not line.strip():
            i += 1
            continue
        
        # Check for multiline format: key<<DELIMITER
        if '<<' in line:
            parts = line.split('<<', 1)
            if len(parts) == 2:
                key = parts[0]
                delimiter = parts[1].strip()
                
                # Read until we find the delimiter
                value_lines = []
                i += 1
                while i < len(lines) and lines[i].strip() != delimiter:
                    value_lines.append(lines[i])
                    i += 1
                
                outputs[key] = '\n'.join(value_lines)
                i += 1  # Skip the closing delimiter
                continue
        
        # Single-line format: key=value
        if '=' in line:
            key, value = line.split('=', 1)
            outputs[key] = value
        
        i += 1
    
    return outputs


def compare_outputs(
    bash_outputs: dict[str, str],
    go_outputs: dict[str, str],
    keys_to_compare: list[str] | None = None,
) -> dict[str, tuple[str, str]]:
    """
    Compare outputs from Bash and Go implementations.
    
    Args:
        bash_outputs: Outputs from Bash implementation
        go_outputs: Outputs from Go implementation
        keys_to_compare: Specific keys to compare (None = all keys)
    
    Returns:
        Dictionary of differing keys with (bash_value, go_value) tuples.
        Empty dict means outputs match.
    """
    differences = {}
    
    if keys_to_compare is None:
        keys_to_compare = list(set(bash_outputs.keys()) | set(go_outputs.keys()))
    
    for key in keys_to_compare:
        bash_val = bash_outputs.get(key, "<missing>")
        go_val = go_outputs.get(key, "<missing>")
        
        if bash_val != go_val:
            differences[key] = (bash_val, go_val)
    
    return differences
