import { RouteGuard } from "../providers/RouteGuard";

export default function HomePage() {

  return (
    <RouteGuard>
      <div>Homepage</div>
    </RouteGuard>
  );
}
