interface AuctionCardProps {
  id: number;
  title: string;
  imageUrl: string;
  bids: number;
}

export default function AuctionCard({
  id,
  title,
  imageUrl,
  bids,
}: AuctionCardProps) {
  return (
    <div className="bg-white shadow-sm rounded-xl overflow-hidden hover:shadow-md transition">
      <img src={imageUrl} alt={title} className="w-full h-48 object-cover" />
      <div className="p-4">
        <h3 className="font-semibold text-lg text-gray-800">{title}</h3>
        <p className="text-sm text-gray-500 mt-1">{bids} bids</p>
        <button className="mt-3 text-sm text-blue-600 hover:underline">
          View Listing
        </button>
      </div>
    </div>
  );
}
