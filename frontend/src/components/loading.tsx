import { Loader } from "lucide-react";

const Loading = () => {
  return (
    <div className="flex items-center justify-center">
      <Loader className="animate-spin w-10 h-10" /> Loading...
    </div>
  );
};

export default Loading;
