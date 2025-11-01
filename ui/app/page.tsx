import Image from "next/image";
import { Button } from "@/components/ui/button";
import Link from "next/link";

export default function Home() {
  return (
    <div className="flex justify-center">
      <div className="flex flex-col justify-center mt-10">
        <h1 className="text-2xl">Twitter Clone</h1>
        <Button asChild>
          <Link href="/login">Login</Link>
        </Button>
      </div>
    </div>
  );
}
