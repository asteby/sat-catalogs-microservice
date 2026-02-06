#!/usr/bin/env python3
import os
import glob

def robust_split_values(content):
    """Extrae las filas (values) respetando paréntesis internos y comillas."""
    values = []
    current_row = []
    in_string = False
    parenthesis_level = 0
    
    # Buscamos donde empiezan los VALUES
    start_idx = content.upper().find("VALUES")
    if start_idx == -1: return []
    
    # Empezamos a leer después de "VALUES"
    data = content[start_idx + 6:].strip()
    
    temp = ""
    for char in data:
        if char == "'" and (not temp or temp[-1] != "\\"): # Manejo de comillas
            in_string = not in_string
        
        if not in_string:
            if char == "(":
                parenthesis_level += 1
            elif char == ")":
                parenthesis_level -= 1
                if parenthesis_level == 0:
                    temp += char
                    values.append(temp.strip())
                    temp = ""
                    continue
        
        if parenthesis_level > 0:
            temp += char
            
    return values

def process_file(filepath, chunk_size=1000):
    print(f"Procesando de forma robusta: {filepath}...")
    
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()

    # Extraer nombre de la tabla
    import re
    table_match = re.search(r'INSERT\s+INTO\s+(\w+)', content, re.I)
    if not table_match: return
    table_name = table_match.group(1)

    # Extraer valores con lógica de paréntesis
    all_rows = robust_split_values(content)
    
    if not all_rows:
        print(f"  --> No se encontraron filas en {filepath}")
        return

    # Escribir con chunking
    with open(filepath, 'w', encoding='utf-8') as f:
        for i in range(0, len(all_rows), chunk_size):
            chunk = all_rows[i : i + chunk_size]
            f.write(f"INSERT INTO {table_name} VALUES\n")
            # Unir filas, asegurando que cada una termine en coma excepto la última del bloque
            clean_chunk = [row.rstrip(',').rstrip(';') for row in chunk]
            f.write(",\n".join(clean_chunk))
            f.write(";\n\n")

    print(f"  --> Éxito: {len(all_rows)} filas procesadas.")

def main():
    data_dir = 'database/data'
    sql_files = glob.glob(os.path.join(data_dir, '*.sql'))
    for filepath in sql_files:
        process_file(filepath)

if __name__ == '__main__':
    main()