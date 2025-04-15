import express, { NextFunction, Request, Response } from "express";
import jwt from "jsonwebtoken";
import bcyrpt from "bcryptjs";
import { getUserByEmail, createUser, User } from "../models/user";
import { validateRequest } from "../middleware/validateRequest";
import { authSchema } from "../schemas/authSchemas";
import { BadRequestError } from "../middleware/errors/bad-request-error";

const router = express.Router();

router.post(
  "/api/users/signup",
  validateRequest(authSchema),
  async (req: Request, res: Response, next: NextFunction) => {
    const { email, password } = req.body;

    // Check if user exists already
    const existingUser = await getUserByEmail(email);
    if (existingUser) {
      return next(new BadRequestError("Email in use"));
    }

    const passwordHash = await bcyrpt.hash(password, 10);

    const user: User = { email, passwordHash, role: "user" };
    await createUser(user);

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

export { router as signupRouter };
