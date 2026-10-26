# <Project name>: Design

> Fill this in before writing feature code. It is the spec every feature is built from.

## Assumptions / clarifications
- 

## Entities & relationships
- 

## State machine (if any)
```
A → B → C
B → D
Anything else → 409 INVALID_TRANSITION
```

## Business rules (numbered: reference these in prompts)
1. 
2. 

## API contracts
| Method | Path | Success | Errors |
|---|---|---|---|
|  |  |  |  |

Error format: `{"error":{"code":"...","message":"...","details":[{"field":"...","issue":"..."}]}}`
List format:  `{"data":[...],"next_cursor":"...","limit":20}`

## Schema
- Tables / columns:
- Constraints (NOT NULL, CHECK, UNIQUE, FK):
- Indexes (and the query each one serves):

## Concurrency / consistency
- Race risks and how they're handled:

## Build order (vertical slices)
1. 
2. 

## Limitations / next steps
- 
