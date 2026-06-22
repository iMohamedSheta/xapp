import main from "./main";
import settings from './settings';
import auth from './auth';
import dashboard from './dashboard';
import wallet from './wallet';
import common from "./common";
import invoices from "./invoices";
import permissions from "./permissions";
import users from "./users";
import team from './team';
import sidebar from './sidebar';
export default {
  sidebar,
  team,
  users,
  wallet,
  main,
  settings,
  auth,
  dashboard,
  common,
  invoices,
  permissions,
} as const
