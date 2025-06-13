"use client";

import { BidButton } from "./BidButton";

interface ListingDetailProps {
  listing: {
    id: string;
    title: string;
    subtitle?: string;
    description: string;
    imageUrl: string;
    specs: Record<string, string>;
    endTime: string;
    seller: {
      name: string;
      location: string;
    };
  };
}

export default function AuctionDetail({ listing }: ListingDetailProps) {
  return (
    <div className="max-w-5xl mx-auto p-6 space-y-8">
      {/* Image */}
      <div className="rounded-2xl overflow-hidden shadow-lg">
        <img
          src={listing.imageUrl}
          alt={listing.title}
          className="w-full h-[450px] object-cover"
        />
      </div>

      {/* Title + Action */}
      <div className="flex justify-between items-start">
        <div>
          <h1 className="text-3xl font-semibold">{listing.title}</h1>
          {listing.subtitle && (
            <p className="text-gray-500">{listing.subtitle}</p>
          )}
          <p className="text-sm mt-1 text-muted-foreground">
            Auction ends: {new Date(listing.endTime).toLocaleString()}
          </p>
        </div>
        <BidButton />
      </div>

      {/* Specs */}
      <div className="grid grid-cols-2 md:grid-cols-3 gap-4 bg-gray-50 p-4 rounded-xl shadow-sm">
        {Object.entries(listing.specs).map(([label, value]) => (
          <div key={label}>
            <p className="text-sm text-muted-foreground">{label}</p>
            <p className="text-base font-medium">{value}</p>
          </div>
        ))}
      </div>

      {/* Seller */}
      <div className="p-4 rounded-xl border">
        <h2 className="text-lg font-semibold mb-1">Seller</h2>
        <p>{listing.seller.name}</p>
        <p className="text-sm text-muted-foreground">
          {listing.seller.location}
        </p>
      </div>

      {/* Description */}
      <div>
        <h2 className="text-xl font-semibold mb-2">Description</h2>
        <p className="text-gray-800 leading-relaxed whitespace-pre-line">
          {listing.description}
        </p>
      </div>
    </div>
  );
}
