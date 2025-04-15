interface FeaturedAuctionProps {
  title: string;
  imageUrl: string;
  description: string;
}

export default function FeaturedAuction({
  title,
  imageUrl,
  description,
}: FeaturedAuctionProps) {
  return (
    <div className="relative bg-white rounded-2xl shadow-md overflow-hidden">
      <img src={imageUrl} alt={title} className="w-full h-72 object-cover" />
      <div className="p-6">
        <h1 className="text-3xl font-bold text-gray-900">{title}</h1>
        <p className="text-gray-600 mt-2">{description}</p>
        <button className="mt-4 px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition">
          View Auction
        </button>
      </div>
    </div>
  );
}
