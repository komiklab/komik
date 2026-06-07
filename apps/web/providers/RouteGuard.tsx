"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useGetAdmin } from "../api/komik";

export function RouteGuard({ children }: { children: React.ReactNode }) {
    const { data, isLoading, isError } = useGetAdmin();
    const router = useRouter();

    useEffect(() => {
        if (isLoading) return;

        if (isError) {
            router.replace("/error");
            return;
        }

        if (data?.status === 200 && !data.data.exists) {
            router.replace("/bootstrap");
        }
    }, [data, isLoading, isError, router]);

    if (isLoading || (data?.status === 200 && !data.data.exists)) {
        return <div>Loading...</div>;
    }

    return <>{children}</>;
}
