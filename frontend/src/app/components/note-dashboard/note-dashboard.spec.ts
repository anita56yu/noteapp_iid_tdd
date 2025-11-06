import { ComponentFixture, TestBed } from '@angular/core/testing';

import { NoteDashboard } from './note-dashboard';

describe('NoteDashboard', () => {
  let component: NoteDashboard;
  let fixture: ComponentFixture<NoteDashboard>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [NoteDashboard]
    })
    .compileComponents();

    fixture = TestBed.createComponent(NoteDashboard);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
