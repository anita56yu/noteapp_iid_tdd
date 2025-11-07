import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { NoteDashboard } from './components/note-dashboard/note-dashboard';
import { NoteViewComponent } from './components/note-view/note-view';

const routes: Routes = [
  { path: '', redirectTo: '/dashboard', pathMatch: 'full' },
  { path: 'dashboard', component: NoteDashboard },
  { path: 'notes/:id', component: NoteViewComponent },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
