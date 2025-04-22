import express, { Request, Response, NextFunction } from "express";
import { authenticateUser } from "../middleware/authenticateUser";
import axios from "axios";
import { BadRequestError } from "../middleware/errors/bad-request-error";

const router = express.Router();

router.use(authenticateUser);

router.use(
  "/api/listings/*",
  async (req: Request, res: Response, next: NextFunction) => {
    try {
      const listingsAddr = process.env.LISTINGS_SERVICE || "";
      const listingsPort = process.env.LISTING_SERVICES_PORT || 3000;

      const method = req.method.toLowerCase();

      // TODO: include headers for now. Will change
      // in the future things like logging...
      const headers = { ...req.headers, host: "listings-srv" };
      const data = req.body;

      const url = `http://${listingsAddr}:${listingsPort}${req.originalUrl}`;

      const response = await axios.request({
        method,
        url,
        headers,
        data,
      });

      res.status(response.status).send(response.data);
    } catch (err) {
      console.error("Error proxying to listings service:", err);
      throw new BadRequestError("Internal Error...");
    }
  }
);

export { router as serviceRouter };
