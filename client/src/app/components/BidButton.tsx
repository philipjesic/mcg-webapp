"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Button } from "./Button";

export const BidButton = () => {
  const router = useRouter();
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [bidAmount, setBidAmount] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleOpen = () => {
    // TODO: check cookies for user to be signed
    // in before making a bid
    /*const token = Cookies.get("session");

    if (!token) {
      alert("Please sign in to place a bid.");
      router.push("/signin");
    } else {
      setIsModalOpen(true);
    }*/
    setIsModalOpen(true);
  };

  const handleSubmit = async () => {
    const amount = parseFloat(bidAmount);
    if (isNaN(amount) || amount <= 0) {
      alert("Please enter a valid bid amount.");
      return;
    }

    try {
      setIsSubmitting(true);
      // Replace with real bid API call
      // await axios.post("/api/bid", { listingId, amount });

      alert("✅ Bid submitted successfully!");
      setIsModalOpen(false);
      setBidAmount("");
    } catch (err) {
      console.error(err);
      alert("❌ Failed to submit bid.");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <>
      <Button onClick={handleOpen} size="lg">
        Place Bid
      </Button>

      {isModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center">
          {/* Backdrop overlay with blur and dark tint */}
          <div className="absolute inset-0 bg-black/40 backdrop-blur-sm" />

          {/* Modal content */}
          <div className="relative z-10 bg-white rounded-lg p-6 w-full max-w-md shadow-xl">
            <h2 className="text-xl font-bold mb-4">Enter Your Bid</h2>

            <input
              type="number"
              value={bidAmount}
              onChange={(e) => setBidAmount(e.target.value)}
              placeholder="Amount in USD"
              className="w-full px-4 py-2 border rounded mb-4"
            />

            <div className="flex justify-end gap-2">
              <Button
                variant="secondary"
                onClick={() => {
                  setIsModalOpen(false);
                  setBidAmount("");
                }}
              >
                Cancel
              </Button>
              <Button onClick={handleSubmit} disabled={isSubmitting}>
                {isSubmitting ? (
                  <div className="flex items-center gap-2">
                    <span className="h-4 w-4 border-2 border-t-2 border-gray-200 rounded-full animate-spin" />
                    Submitting...
                  </div>
                ) : (
                  "Submit Bid"
                )}
              </Button>
            </div>
          </div>
        </div>
      )}
    </>
  );
};
