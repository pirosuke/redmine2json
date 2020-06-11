# redmine2json
Fetch Redmine Tickets And Save As JSON File.

## Usage

```
rm2json [flags]
  -c string
        Config Name
  -o string
        Output Directory Path (default ".")
```

## Config File Sample

```json
{
    "redmine": {
        "api_key": "REDMINE API KEY",
        "url_root": "https://redmine/url/root",
        "basic_auth": {
            "username": "",
            "password": ""
        }
    }
}
```
