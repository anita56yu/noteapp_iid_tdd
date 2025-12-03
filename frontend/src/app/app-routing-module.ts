import { Routes } from '@angular/router';
import { NoteDashboard } from './components/note-dashboard/note-dashboard';
import { LoginPage } from './components/login-page/login-page';
import { AuthGuard } from './auth.guard';

export const routes: Routes = [
  { path: 'login', component: LoginPage },
  { path: '', component: NoteDashboard, canActivate: [AuthGuard] },
  { path: '**', redirectTo: '' }
];
