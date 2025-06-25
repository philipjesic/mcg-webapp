"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Button } from "./Button";
import useRequest from "../hooks/use-request";
import { BidRequest } from "../requests/bids";

export const BidButton = ({ listingId }: { listingId: string }) => {
  const router = useRouter();
  const [modalStep, setModalStep] = useState<"entry" | "confirm" | "success">(
    "entry"
  );
  const [isAnimatingOut, setIsAnimatingOut] = useState(false);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [bidAmount, setBidAmount] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  const bidRequest: BidRequest = {
    data: {
      auction_id: listingId,
      user_id: "current_user_id", // Replace with actual user ID logic
      amount: parseFloat(bidAmount),
      timestamp: new Date().toISOString(),
    },
  };

  const { doRequest, errors } = useRequest({
    url: `/api/bids`,
    method: "post",
    body: bidRequest,
    onSuccess: () => {
      setModalStep("success");
    },
  });

  const handleSubmit = async () => {
    const amount = parseFloat(bidAmount);
    if (isNaN(amount) || amount <= 0) {
      alert("Please enter a valid bid amount.");
      return;
    }

    setIsSubmitting(true);
    await doRequest();
    setIsSubmitting(false);
  };

  const closeModal = () => {
    setIsAnimatingOut(true);
    setTimeout(() => {
      setIsAnimatingOut(false);
      setIsModalOpen(false);
      setModalStep("entry");
      setBidAmount("");
    }, 400); // match the transition duration
  };

  return (
    <>
      <Button onClick={() => setIsModalOpen(true)} size="lg">
        Place Bid
      </Button>

      {isModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center">
          <div className="absolute inset-0 bg-black/40 backdrop-blur-sm" />

          {/* Animated modal wrapper */}
          <div
            className={`relative z-10 bg-white rounded-lg p-6 w-full max-w-md shadow-xl transform transition-all duration-150 ease-out
        ${isAnimatingOut ? "opacity-0 scale-95" : "opacity-100 scale-100"}
      `}
          >
            {modalStep === "entry" && (
              <>
                <h2 className="text-xl font-bold mb-4">Enter Your Bid</h2>
                <input
                  type="number"
                  value={bidAmount}
                  onChange={(e) => setBidAmount(e.target.value)}
                  placeholder="Amount in USD"
                  className="w-full px-4 py-2 border rounded mb-4"
                />
                <div className="flex justify-end gap-2">
                  <Button variant="secondary" onClick={closeModal}>
                    Cancel
                  </Button>
                  <Button
                    onClick={() => {
                      const amt = parseFloat(bidAmount);
                      if (isNaN(amt)) {
                        alert("Please enter a valid bid amount.");
                      } else {
                        setModalStep("confirm");
                      }
                    }}
                  >
                    Next
                  </Button>
                </div>
              </>
            )}

            {modalStep === "confirm" && (
              <>
                <h2 className="text-xl font-bold mb-2">Confirm Your Bid</h2>
                <p className="mb-2">
                  You're placing a bid of <strong>${bidAmount}</strong> on:
                </p>
                <p className="italic text-gray-700 mb-4">"{listingId}"</p>
                {errors && <div className="mb-2">{errors}</div>}
                <div className="flex justify-end gap-2">
                  <Button variant="secondary" onClick={closeModal}>
                    Cancel
                  </Button>
                  <Button onClick={handleSubmit} disabled={isSubmitting}>
                    {isSubmitting ? (
                      <div className="flex items-center gap-2">
                        <span className="h-4 w-4 border-2 border-t-2 border-gray-200 rounded-full animate-spin" />
                        Submitting...
                      </div>
                    ) : (
                      "Confirm Bid"
                    )}
                  </Button>
                </div>
              </>
            )}

            {modalStep === "success" && (
              <>
                <h2 className="text-xl font-bold mb-4 text-green-600">
                  ✅ Bid Submitted
                </h2>
                <p className="mb-4">
                  Your bid of <strong>${bidAmount}</strong> for{" "}
                  <em>{listingId}</em> has been placed successfully.
                </p>
                <div className="flex justify-end">
                  <Button onClick={closeModal}>Close</Button>
                </div>
              </>
            )}
          </div>
        </div>
      )}
    </>
  );
};
