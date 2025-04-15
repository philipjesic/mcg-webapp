import { Request, Response, NextFunction } from "express";
import jwt from "jsonwebtoken";
import { NotAuthorizedError } from "./errors/not-authorized-error";

const authenticateUser = (req: Request, res: Response, next: NextFunction) => {
  const token = req.cookies.jwt;
  if (!token) {
    throw new NotAuthorizedError();
  }
  try {
    jwt.verify(token, process.env.JWT_SECRET as string);
    next();
  } catch (error) {
    throw new NotAuthorizedError();
  }
};

export { authenticateUser };
