import { Component, signal } from '@angular/core';
import { RouterModule } from '@angular/router';
import { HttpClientModule } from '@angular/common/http';
import { NoteDashboard } from './components/note-dashboard/note-dashboard';

@Component({
  selector: 'app-root',
  templateUrl: './app.html',
  standalone: true,
  imports: [NoteDashboard, RouterModule, HttpClientModule],
  styleUrl: './app.scss'
})
export class App {
  protected readonly title = signal('noteapp-frontend');
}
