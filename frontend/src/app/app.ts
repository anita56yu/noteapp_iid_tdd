import { Component, signal } from '@angular/core';
import { RouterModule } from '@angular/router';
import { HttpClientModule } from '@angular/common/http';

@Component({
  selector: 'app-root',
  templateUrl: './app.html',
  standalone: true,
  imports: [RouterModule, HttpClientModule],
  styleUrl: './app.scss'
})
export class App {
  protected readonly title = signal('noteapp-frontend');
}
