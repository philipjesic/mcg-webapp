import express, { NextFunction, Request, Response } from "express";
import jwt from "jsonwebtoken";
import bcrypt from "bcryptjs";
import { getUserByEmail } from "../models/user";
import { validateRequest } from "../middleware/validateRequest";
import { authSchema } from "../schemas/authSchemas";
import { BadRequestError } from "../middleware/errors/bad-request-error";

const router = express.Router();

router.post(
  "/api/users/login",
  validateRequest(authSchema),
  async (req: Request, res: Response, next: NextFunction) => {
    const { email, password } = req.body;

    const existingUser = await getUserByEmail(email);

    if (!existingUser) {
      return next(new BadRequestError("Invalid credentials"));
    }

    const passwordsMatch = await bcrypt.compare(
      password,
      existingUser.passwordHash
    );

    if (!passwordsMatch) {
      return next(new BadRequestError("Invalid credentials"));
    }

    // Generate JWT
    const tokenClaims = { email, role: "user" };
    const token = jwt.sign(tokenClaims, process.env.JWT_SECRET as string, {
      expiresIn: "15m",
    });

    // Set cookie
    res.cookie("jwt", token, {
      httpOnly: true,
    });

    res.status(201).json({ data: tokenClaims });
  }
);

export { router as loginRouter };
