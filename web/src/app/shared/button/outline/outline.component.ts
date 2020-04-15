import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';
import { ButtonComponent } from '../button';

@Component({
  selector: 'app-outline-button',
  changeDetection: ChangeDetectionStrategy.OnPush,
  templateUrl: './outline.component.html',
  styleUrls: ['./outline.component.scss'],
})
export class OutlineButtonComponent extends ButtonComponent { }
