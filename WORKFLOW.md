# Workflow

Volg deze workflow alsjeblieft om een gestructureerde codebase te houden. 

### 1. Lees documentatie in alle packages (README.md)
- In elke package staat een README.md, lees deze door om er achter te komen, wat het doet, waarvoor het is én hoe het werkt.

### 2. Schrijf je code
- Voeg toe aan het project door je taken uit te voeren
- Gebruik zoveel mogelijk het principe DRY (Don't repeat yourself), gebruik bestaande functies zoveel mogelijk, deze functies ken je door het lezen van de documentatie:)
- Gebruik vaak git pull, zo voorkomen we problemen met merging.
- Maak een nieuwe branch aan met jouw veranderingen óf gebruik een bestaande branch:
```bash
# Zo maak je een nieuwe branch aan
git pull
git checkout -b feature/<feature_naam>
# Of gebruik een bestaande branch
git checkout <branch_naam>
git pull
# Voeg je veranderingen toe (Doe dit niet in /apps/server maar gewoon in de root van het project)
git add .
# Commit je veranderingen
git commit -m "hier vertel je wat je hebt gedaan"
# Push naar GitHub
git push -u origin <branch_naam>
```
- Heb je vragen? Gebruik dan de groepsapp of gebruik AI in tijdsnood, met voorkeur aan Claude 4.

### 3. Klaar met een feature?
- Maak een PR aan (Pull request) in GitHub.
- Vertel het aan de anderen want tenminste één iemand moet het goedkeuren. Zo controleren we elkaar en houden we de codebase netjes.
- Merge!
