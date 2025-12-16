import FeaturedAuction from "./components/FeaturedAuction";
import AuctionCard from "./components/AuctionCard";
import { Listing, ListingsResponse } from "./responses/listings";

async function getListings(): Promise<Listing[]> {
  const listingsAddr = process.env.LISTINGS_SERVICE || "";
  const listingsPort = process.env.LISTING_SERVICES_PORT || 3000;
  const res = await fetch(
    `http://${listingsAddr}:${listingsPort}/api/listings`
  );
  if (!res.ok) {
    console.log("oh no error...");
  }
  const response: ListingsResponse = await res.json()
  return response.data
}

export default async function Home() {
  const listings = await getListings();
  return (
    <div className="px-4 py-8 max-w-7xl mx-auto">
      <FeaturedAuction
        title="2023 Porsche 911 GT3 Touring"
        imageUrl="/porsche.jpg"
        description="Rare 6-speed manual in Shark Blue – ends in 3 hours!"
      />

      <h2 className="text-2xl font-semibold mt-12 mb-6 text-gray-800">
        Ongoing Auctions
      </h2>

      <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-6">
        {listings.map((listing) => (
          <AuctionCard key={listing.id} {...listing} />
        ))}
      </div>
    </div>
  );
}
