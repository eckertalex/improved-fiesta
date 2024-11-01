# Improved Fiesta

| Method   | URL Pattern                 | Action                                          |
| -------- | --------------------------- | ----------------------------------------------- |
| `GET`    | `/v1/healthcheck`           | Show application health and version information |
| `POST`   | `/v1/users`                 | Register a new user                             |
| `PUT`    | `/v1/users/activated`       | Activate a specific user                        |
| `PUT`    | `/v1/users/password`        | Update the password for a specific user         |
| `POST`   | `/v1/tokens/authentication` | Generate a new authentication token             |
| `POST`   | `/v1/tokens/password-reset` | Generate a new password-reset token             |
| `GET`    | `/debug/vars`               | Display application metrics                     |
