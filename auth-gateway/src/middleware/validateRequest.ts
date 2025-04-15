import { Request, Response, NextFunction } from "express";
import { ZodError, ZodSchema } from "zod";
import { RequestValidationError } from "./errors/request-validation-error";

export const validateRequest = (schema: ZodSchema) => {
  return (req: Request, res: Response, next: NextFunction): void => {
    try {
      schema.parse(req.body);
      next();
    } catch (error) {
      if (error instanceof ZodError) {
        throw new RequestValidationError(error.issues);
      }
      res.status(500).json({ message: "Internal Server Error" });
      return; // Explicit return
    }
  };
};
