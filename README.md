# Improved Fiesta

| Method  | URL Pattern                 | Action                                          | Body                                         |
| ------- | --------------------------- | ----------------------------------------------- | -------------------------------------------- |
| `GET`   | `/v1/healthcheck`           | Show application health and version information |                                              |
| `POST`  | `/v1/users`                 | Register a new user                             | `{name:string,email:string,password:string}` |
| `PUT`   | `/v1/users/activated`       | Activate a specific user                        | `{token:string}`                             |
| `PUT`   | `/v1/users/password`        | Update the password for a specific user         | `{password:string,token:string}`             |
| `PATCH` | `/v1/users/:id/role`        | Update the role for a specific user             | `{role:'admin'/'user'}`                     |
| `POST`  | `/v1/tokens/authentication` | Generate a new authentication token             | `{email:string,password:string}`             |
| `POST`  | `/v1/tokens/activation`     | Generate a new activation token                 | `{email:string}`                             |
| `POST`  | `/v1/tokens/password-reset` | Generate a new password-reset token             | `{email:string}`                             |
| `GET`   | `/debug/vars`               | Display application metrics                     |                                              |
