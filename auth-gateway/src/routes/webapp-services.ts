import express, { Request, Response, NextFunction } from "express";
import { authenticateUser } from "../middleware/authenticateUser";
import axios from "axios";
import { BadRequestError } from "../middleware/errors/bad-request-error";
import { ServerError } from "../middleware/errors/server-error";

const router = express.Router();

router.use(authenticateUser);

router.use(
  "/api/listings*",
  async (req: Request, res: Response, next: NextFunction) => {
    console.log("got to handler");
    try {
      const listingsAddr = process.env.LISTINGS_SERVICE || "";
      const listingsPort = process.env.LISTING_SERVICE_PORT || 3000;

      const method = req.method.toLowerCase();

      // TODO: include headers for now. Will change
      // in the future things like logging...
      //const headers = { ...req.headers, host: "listings-srv" };
      const data = req.body;

      const url = `http://${listingsAddr}:${listingsPort}${req.originalUrl}`;

      const response = await axios.request({
        method,
        url,
        //headers,
        data,
      });
      console.log("resolved response");

      res.status(response.status).send(response.data);
    } catch (err) {
      console.error("Error proxying to listings service:", err);
      return next(new ServerError("Internal Error..."));
    }
  }
);

router.use(
  "/api/bids*",
  async (req: Request, res: Response, next: NextFunction) => {
    console.log("got to handler");
    try {
      const bidsAddr = process.env.BIDS_SERVICE || "";
      const bidsPort = process.env.BIDS_SERVICE_PORT || 3000;

      const method = req.method.toLowerCase();

      // TODO: include headers for now. Will change
      // in the future things like logging...
      //const headers = { ...req.headers, host: "listings-srv" };
      const data = req.body;

      const url = `http://${bidsAddr}:${bidsPort}${req.originalUrl}`;

      const response = await axios.request({
        method,
        url,
        //headers,
        data,
      });

      res.status(response.status).send(response.data);
    } catch (err) {
      console.error("Error proxying to bids service:", err);
      return next(new ServerError("Internal Error..."));
    }
  }
);

export { router as serviceRouter };
