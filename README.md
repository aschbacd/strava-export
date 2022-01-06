# strava-export

Excel exporter for Strava.

## Configuration

The application can be configured using environment variables, either by using the file `.env` to
store key/value pairs or by directly exporting them in the applications environment. The following
variables can be used:

| Environment variable | Description                                       | Default                 |
| -------------------- | ------------------------------------------------- | ----------------------- |
| ADDRESS              | Address used to launch server                     | `localhost`             |
| PORT                 | Port used to launch server                        | `8080`                  |
| DEBUG                | Enable debug logging for http server              | `false`                 |
| STRAVA_CLIENT_ID     | Strava Application client id                      | `-`                     |
| STRAVA_CLIENT_SECRET | Strava Application client secret                  | `-`                     |
| BASE_URL             | Base url for application (used for auth redirect) | `http://localhost:8080` |

## Swagger client library

Strava provides a swagger spec to generate client libraries for their api. The following command
was used to generate to go library and store it in a package:

```bash
swagger-codegen generate -i https://developers.strava.com/swagger/swagger.json -l go -o pkg/strava
```
