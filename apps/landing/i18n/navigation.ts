// SPDX-License-Identifier: MIT OR Apache-2.0
import { createNavigation } from "next-intl/navigation";
import { routing } from "./routing";

export const { Link, usePathname } = createNavigation(routing);
