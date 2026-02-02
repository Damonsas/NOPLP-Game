import json
import re
import os

def texte_en_json(texte):
    # Extraire titre et artiste
    titre_match = re.search(r'^titre\s*:\s*(.+)$', texte, flags=re.IGNORECASE | re.MULTILINE)
    artiste_match = re.search(r'^artiste\s*:\s*(.+)$', texte, flags=re.IGNORECASE | re.MULTILINE)
    titre = titre_match.group(1).strip() if titre_match else "Inconnu"
    artiste = artiste_match.group(1).strip() if artiste_match else "Inconnu"

    # Séparer les sections (couplet/refrain), insensible à la casse et aux points
    pattern = re.compile(r'^(couplet\d*|refrain\d*)[\.]?', re.IGNORECASE | re.MULTILINE)
    splits = pattern.split(texte)
    
    paroles = {}
    i = 1
    while i < len(splits):
        key = splits[i].lower().replace('.', '').replace(' ', '')
        content = splits[i+1].strip()
        lines = [line.strip() for line in content.split('\n') if line.strip()]
        paroles[key] = lines
        i += 2

    return {
        "titre": titre,
        "artiste": artiste,
        "paroles": paroles
    }

    ## format du truc : couplet1 sans : 

# Nom du fichier à convertir
nom_fichier_txt = os.path.join(os.path.dirname(__file__), "Clair Obscur - Lumière.txt")
if not os.path.exists(nom_fichier_txt):
    print(f"Le fichier {nom_fichier_txt} n'existe pas dans ce dossier.")
else:
    with open(nom_fichier_txt, "r", encoding="utf-8") as f:
        texte = f.read()

    json_data = texte_en_json(texte)

    nom_fichier_json = os.path.join(os.path.dirname(__file__), "Clair Obscur - Lumière.json")
    with open(nom_fichier_json, "w", encoding="utf-8") as f:
        json.dump(json_data, f, ensure_ascii=False, indent=2)

    print(f"JSON créé avec succès : {nom_fichier_json}")


