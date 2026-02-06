#!/usr/bin/env python3
import os
import re
import glob

def batch_inserts_in_file(filepath):
    print(f"Leyendo: {filepath}...")
    
    with open(filepath, 'r', encoding='utf-8') as f:
        lines = f.readlines()

    inserts_by_table = {}
    other_lines = []

    # Regex mejorado: 
    # - Ignora mayúsculas/minúsculas (re.I)
    # - Soporta espacios variables entre INTO, Tabla y VALUES
    # - Captura todo lo que esté dentro de los paréntesis de VALUES
    pattern = re.compile(r'INSERT\s+INTO\s+(\w+)\s+VALUES\s*\((.*)\);', re.I)

    for line in lines:
        line = line.strip()
        if not line:
            continue
            
        match = pattern.match(line)
        if match:
            table = match.group(1)
            values = match.group(2)
            if table not in inserts_by_table:
                inserts_by_table[table] = []
            inserts_by_table[table].append(f"({values})")
        else:
            # Guardamos líneas que no son INSERT (CREATE TABLE, comentarios, etc.)
            if line:
                other_lines.append(line)

    if not inserts_by_table:
        print(f"  --> No se encontraron INSERTs válidos en {filepath}")
        return

    # Escribir el resultado
    with open(filepath, 'w', encoding='utf-8') as f:
        # Primero ponemos las líneas que no eran inserts (estructuras)
        if other_lines:
            f.write('\n'.join(other_lines) + '\n\n')
        
        # Luego los inserts agrupados
        for table, value_list in inserts_by_table.items():
            # Un solo INSERT con todos los VALUES separados por coma
            f.write(f"INSERT INTO {table} VALUES \n" + ",\n".join(value_list) + ";\n")

    print(f"  --> ¡Éxito! {len(inserts_by_table)} tablas optimizadas.")

def main():
    # Ajusta esta ruta a tu estructura de carpetas de Asteby
    data_dir = 'database/data'
    sql_files = glob.glob(os.path.join(data_dir, '*.sql'))
    
    if not sql_files:
        print(f"No se encontraron archivos .sql en {data_dir}")
        return

    for filepath in sql_files:
        batch_inserts_in_file(filepath)

if __name__ == '__main__':
    main()