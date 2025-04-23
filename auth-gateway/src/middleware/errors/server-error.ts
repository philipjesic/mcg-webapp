import { APIError } from "./api-error";

export class ServerError extends APIError {
  statusCode = 500;

  constructor(message: string) {
    super(message);
    Object.setPrototypeOf(this, ServerError.prototype);
  }

  serializeErrors(): { message: string }[] {
    return [{ message: this.message }];
  }
}
