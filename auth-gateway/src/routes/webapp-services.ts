import express, { Request, Response, NextFunction } from "express";
import { authenticateUser } from "../middleware/authenticateUser";

const router = express.Router();

router.use(authenticateUser);

export { router as serviceRouter };


