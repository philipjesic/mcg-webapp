import express from "express";
import dotenv from 'dotenv';
import { json } from "body-parser";
import cookieParser from "cookie-parser";
import { signupRouter } from "./routes/signup";
import { loginRouter } from "./routes/login";
import { serviceRouter } from "./routes/webapp-services";
import { errorHandler } from "./middleware/errors/error-handler";

dotenv.config();

const app = express();
app.use(json());
app.use(cookieParser());

app.use(signupRouter);
app.use(loginRouter);
app.use(serviceRouter);

app.use(errorHandler);

export { app };
