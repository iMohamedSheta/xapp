import { NavItem } from "@/types";
import { ActivityIcon, ArchiveIcon, BarChart2, Bell, BookOpen, Box, BoxesIcon, Building, Calendar, ChartAreaIcon, CreditCard, CreditCardIcon, DollarSign, FileArchive, FileText, Folder, Gift, GroupIcon, HandshakeIcon, HomeIcon, IdCardIcon, Layers, LayoutGrid, Map as MapIcon, Network, Package, PackagePlusIcon, PenTool, PhoneIncomingIcon, Puzzle, ReceiptIcon, Server, Settings, ShoppingCart, StickyNote, TicketCheckIcon, TicketsIcon, UserPlus, Users2, Users2Icon, Wallet, Wifi, Zap } from "lucide-vue-next";
import { lang } from "./lang";


export const adminSidebarMenu = (): NavItem[] => [
  {
    title: lang('sidebar.dashboard'),
    href: '/dashboard',
    activeKey: "/dashboard",
    icon: LayoutGrid,
    permission: 'dashboard.view',
  },
  {
    title: lang('sidebar.settings'),
    href: '/settings/profile',
    activeKey: '/settings',
    icon: Settings,
  },
];

export const clientSidebarMenu = (): NavItem[] => [
  {
    title: lang('sidebar.dashboard'),
    href: '/client/dashboard',
    activeKey: "/client/dashboard",
    icon: LayoutGrid,
  },
  {
    title: lang('sidebar.settings'),
    href: '/settings/profile',
    activeKey: '/settings',
    icon: Settings,
  },
];

export const managerSidebarMenu = (): NavItem[] => [
  {
    title: lang('sidebar.dashboard'),
    href: '/manager/dashboard',
    activeKey: "/dashboard",
    icon: LayoutGrid,
    permission: 'dashboard.view',
  },
  {
    title: lang('sidebar.settings'),
    href: '/settings/profile',
    activeKey: '/settings',
    icon: Settings,
  },
];

export const adminFooterMenu = (): NavItem[] => [
  // {
  //   title: 'الدليل',
  //   href: 'https://laravel.com/docs/starter-kits#vue',
  //   icon: BookOpen,
  // },
];

export const clientFooterMenu = (): NavItem[] => [];
export const managerFooterMenu = (): NavItem[] => [];

export function hasPermission(permissions: Record<string, boolean> | undefined, permissionKey: string | undefined): boolean {
  if (!permissionKey) return true;
  if (!permissions) return false;

  // Check for wildcard permission (super admin)
  if (permissions['*']) {
    return true;
  }

  // Check exact permission
  if (permissions[permissionKey]) {
    return true;
  }

  // Check wildcard permissions (e.g., "admin.*" matches "admin.users.create")
  const parts = permissionKey.split('.');
  for (let i = parts.length; i > 0; i--) {
    const wildcardKey = parts.slice(0, i).join('.') + '.*';
    if (permissions[wildcardKey]) {
      return true;
    }
  }

  return false;
}

function filterMenuItems(items: NavItem[], permissions: Record<string, boolean> | undefined, role?: string): NavItem[] {
  function itemVisibleForRole(item: NavItem, role?: string): boolean {
    const itemRole = (item as any).role;
    if (!itemRole) return true;
    if (!role) return false;
    if (typeof itemRole === 'string') return itemRole === role;
    if (Array.isArray(itemRole)) return itemRole.includes(role);
    return true;
  }

  return items
    .filter(item => hasPermission(permissions, item.permission) && itemVisibleForRole(item, role))
    .map(item => {
      if (item.children) {
        return {
          ...item,
          children: filterMenuItems(item.children, permissions, role)
        };
      }
      return item;
    })
    .filter(item => !item.children || item.children.length > 0 || !item.href === undefined); // keep if has children (that are not filtered out) or is a direct link
}

export function getSidebarMenu(role: string, permissions?: Record<string, boolean>): NavItem[] {
  let menu: NavItem[] = [];
  switch (role) {
    case "super_admin":
    case "admin":
      menu = adminSidebarMenu();
      break;

    case "super_manager":
    case "manager":
      menu = managerSidebarMenu();
      break;

    case 'distributor':
      menu = [];
      menu = [...menu, ...adminSidebarMenu()];
      break;

    case "client":
      menu = clientSidebarMenu();
      break;

    default:
      menu = [];
  }

  return filterMenuItems(menu, permissions, role);
}

export const adminContextMenu = (): NavItem[] => {
  return [];
};

// Context menu is when you right click on any where in the page its not available right now
export function getContextMenu(role: string, permissions?: Record<string, boolean>): NavItem[] {
  let menu: NavItem[] = [];
  switch (role) {
    case "super_admin":
    case "admin":
      menu = adminContextMenu();
      break;

    case "super_manager":
    case "manager":
      menu = managerSidebarMenu();
      break;

    case "client":
      menu = clientSidebarMenu();
      break;

    default:
      menu = [];
  }

  return filterMenuItems(menu, permissions, role);
}

export function getSidebarFooterMenu(role: string, permissions?: Record<string, boolean>): NavItem[] {
  let menu: NavItem[] = [];
  switch (role) {
    case "super_admin":
    case "admin":
      menu = adminFooterMenu();
      break;

    case "super_manager":
    case "manager":
      menu = managerFooterMenu();
      break;

    case "client":
      menu = clientFooterMenu();
      break;

    default:
      menu = [];
  }

  return filterMenuItems(menu, permissions, role);
}

export interface SidebarMenuSettings {
  order: string[];
  hidden: string[];
}

export function normalizeSidebarMenuSettings(
  items: NavItem[],
  settings?: Partial<SidebarMenuSettings> | null,
): SidebarMenuSettings {
  const hrefSet = new Set(items.map((item) => item.href));
  const order = (settings?.order ?? []).filter((href) => hrefSet.has(href));
  const hidden = (settings?.hidden ?? []).filter((href) => hrefSet.has(href));

  return { order, hidden };
}

export function applySidebarMenuSettings(
  items: NavItem[],
  settings?: Partial<SidebarMenuSettings> | null,
): NavItem[] {
  const normalized = normalizeSidebarMenuSettings(items, settings);
  const hiddenSet = new Set(normalized.hidden);
  const itemsByHref = new Map(items.map((item) => [item.href, item]));
  const ordered: NavItem[] = [];

  for (const href of normalized.order) {
    const item = itemsByHref.get(href);
    if (!item || hiddenSet.has(href)) continue;
    ordered.push(item);
    itemsByHref.delete(href);
  }

  for (const item of items) {
    if (!itemsByHref.has(item.href) || hiddenSet.has(item.href)) continue;
    ordered.push(item);
  }

  return ordered;
}