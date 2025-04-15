import { ZodIssue } from "zod";
import { APIError } from "./api-error";

export class RequestValidationError extends APIError {
  statusCode = 400;

  constructor(public errors: ZodIssue[]) {
    super("Invalid request parameters");
    Object.setPrototypeOf(this, RequestValidationError.prototype);
  }

  serializeErrors() {
    return this.errors.map((err) => {
      if (err.path && err.path.length > 0) {
        return { message: err.message, field: err.path[0].toString() };
      }
      return { message: err.message };
    });
  }
}
