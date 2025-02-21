
import { useParams } from "react-router-dom";

const Group = () => {
  const { id } = useParams();
  return (
    <div className="animate-fade-in">
      <h1 className="text-3xl font-semibold mb-8">Group {id}</h1>
    </div>
  );
};

export default Group;
