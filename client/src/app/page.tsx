import FeaturedAuction from "./components/FeaturedAuction";
import AuctionCard from "./components/AuctionCard";

export default function Home() {
  const ongoingAuctions = [
    { id: 1, title: "2020 BMW M2 Competition", imageUrl: "/bmw.jpg", bids: 14 },
    { id: 2, title: "2017 Audi RS3", imageUrl: "/audi.jpg", bids: 20 },
    {
      id: 3,
      title: "2021 Tesla Model 3 Performance",
      imageUrl: "/tesla.jpg",
      bids: 8,
    },
    // Add more...
  ];

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
        {ongoingAuctions.map((auction) => (
          <AuctionCard key={auction.id} {...auction} />
        ))}
      </div>
    </div>
  );
}
