# Improved Fiesta

| Method   | URL Pattern                 | Permissions    | Body                                         | Description                                     |
| -------- | --------------------------- | -------------- | -------------------------------------------- | ----------------------------------------------- |
| `GET`    | `/v1/healthcheck`           | Public         |                                              | Show application health and version information |
| `POST`   | `/v1/users`                 | Admin          | `{name:string,email:string,password:string}` | Register a new user                             |
| `GET`    | `/v1/users/:id`             | Admin or Owner |                                              | Get a user                                      |
| `PATCH`  | `/v1/users/:id`             | Admin or Owner | `{name:string,email:string,password:string}` | Update a user                                   |
| `PATCH`  | `/v1/users/:id/role`        | Admin          | `{role:'admin'/'user'}`                      | Update the role of a user                       |
| `DELETE` | `/v1/users/:id`             | Admin or Owner |                                              | Delete a user                                   |
| `PUT`    | `/v1/users/activated`       | Public         | `{token:string}`                             | Activate a user                                 |
| `PUT`    | `/v1/users/password`        | Public         | `{password:string,token:string}`             | Update the password for a user                  |
| `POST`   | `/v1/tokens/authentication` | Public         | `{email:string,password:string}`             | Generate a new authentication token             |
| `POST`   | `/v1/tokens/activation`     | Public         | `{email:string}`                             | Generate a new activation token                 |
| `POST`   | `/v1/tokens/password-reset` | Public         | `{email:string}`                             | Generate a new password-reset token             |
| `GET`    | `/debug/vars`               |                |                                              | Display application metrics                     |
