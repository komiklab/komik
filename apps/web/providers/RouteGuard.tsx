"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useAdminStore } from "../stores/admin";

export function RouteGuard({ children }: { children: React.ReactNode }) {
    const { doesAdminExist, checkIfAdminExists, error } = useAdminStore();
    const [loading, setLoading] = useState(true);
    const router = useRouter();

    useEffect(() => {
        let mounted = true;

        checkIfAdminExists().finally(() => {
            if (mounted) {
                setLoading(false);
            }
        });

        return () => {
            mounted = false;
        };
    }, [checkIfAdminExists]);

    useEffect(() => {
        if (!loading && error) {
            console.log(error)
            router.replace("/error")
            return
        }
        if (!loading && !doesAdminExist) {
            router.replace("/bootstrap");
        }
    }, [doesAdminExist, loading, router, error]);

    if (loading || !doesAdminExist) {
        return <div>Loading...</div>;
    }

    return <>{children}</>;
}
