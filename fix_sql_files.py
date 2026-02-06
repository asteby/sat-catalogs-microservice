import os
import re
import glob

SCHEMA_DIR = './database/schemas'
DATA_DIR = './database/data'

def fix_schema_file(filepath):
    with open(filepath, 'r') as f:
        content = f.read()

    # Generic fixes
    content = content.replace('"', '`')
    
    # Capitalize TEXT types (not strictly necessary but consistent)
    # Using insensitive replacement for " text " to " TEXT "
    content = re.sub(r'\s+text\s+', ' TEXT ', content, flags=re.IGNORECASE)
    content = re.sub(r'\s+text,', ' TEXT,', content, flags=re.IGNORECASE)
    content = re.sub(r'\s+text$', ' TEXT', content, flags=re.IGNORECASE)  # End of line

    # Identify Primary Keys
    # Pattern: PRIMARY KEY(`id`) or PRIMARY KEY(`col1`, `col2`)
    pk_match = re.search(r'PRIMARY KEY\s*\(([^)]+)\)', content, re.IGNORECASE)
    if pk_match:
        pk_cols_str = pk_match.group(1)
        # Split by comma and clean
        pk_cols = [col.strip().strip('`') for col in pk_cols_str.split(',')]
        
        for col in pk_cols:
            # Replace TEXT definition with VARCHAR(255) for this column
            # Regex: `col` TEXT
            # We look for `col` followed by TEXT types
            pattern = re.compile(rf'`{re.escape(col)}`\s+TEXT', re.IGNORECASE)
            content = pattern.sub(f'`{col}` VARCHAR(255)', content)

    # Write back
    with open(filepath, 'w') as f:
        f.write(content)
    print(f"Fixed schema: {filepath}")

def fix_data_file(filepath):
    with open(filepath, 'r') as f:
        lines = f.readlines()

    new_lines = []
    for line in lines:
        stripped = line.strip()
        if stripped.upper().startswith('PRAGMA'):
            continue
        if stripped.upper() in ('BEGIN TRANSACTION;', 'COMMIT;', 'BEGIN;', 'COMMIT'):
            continue
        
        # Replace double quotes in INSERT statements if any (though example showed none)
        # But some might have INSERT INTO "table"
        if stripped.upper().startswith('INSERT INTO'):
            line = line.replace('"', '`')
        
        new_lines.append(line)

    with open(filepath, 'w') as f:
        f.writelines(new_lines)
    print(f"Fixed data: {filepath}")

def main():
    # Fix Schemas
    schema_files = glob.glob(os.path.join(SCHEMA_DIR, '*.sql'))
    for f in schema_files:
        fix_schema_file(f)

    # Fix Data
    data_files = glob.glob(os.path.join(DATA_DIR, '*.sql'))
    for f in data_files:
        fix_data_file(f)

if __name__ == '__main__':
    main()
