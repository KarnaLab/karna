### Configuration

```javascript

{
"global": {},
"deployments": [
  {
    "src": [string - required],
    "key": [string - optional],
    "file": [string - required],
    "functionName": [string - required],
    "aliases": [map - required]{
    "<some-alias>": [string - required]
    },
    "bucket": [string - optional],
    "prune":[map - optional] {
      "keep": [int - optional],
      "alias": [bool - optional]
    }
  }
]
}
```
