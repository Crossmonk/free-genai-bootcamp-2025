
import { useParams } from "react-router-dom";

const StudyActivity = () => {
  const { id } = useParams();
  return (
    <div className="animate-fade-in">
      <h1 className="text-3xl font-semibold mb-8">Study Activity {id}</h1>
    </div>
  );
};

export default StudyActivity;
