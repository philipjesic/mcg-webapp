import { Request, Response, NextFunction } from "express";

export abstract class APIError extends Error {
  abstract statusCode: number;
  
  constructor(message: string) {
    super(message);
    Object.setPrototypeOf(this, APIError.prototype);
  }

  abstract serializeErrors(): { message: string; field?: string }[];
}
