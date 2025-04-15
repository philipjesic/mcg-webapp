import { APIError } from "./api-error";

export class NotAuthorizedError extends APIError {
  statusCode = 401;

  constructor() {
    super("Not Authorized");
    Object.setPrototypeOf(this, NotAuthorizedError.prototype);
  }

  serializeErrors(): { message: string; field?: string }[] {
    return [{ message: "Not Authorized" }];
  }
}
