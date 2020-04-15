import { ChangeDetectionStrategy, Component } from '@angular/core';
import { AbstractPostComponent } from '../abstract-post.component';

@Component({
  selector: 'app-post-tags',
  changeDetection: ChangeDetectionStrategy.OnPush,
  templateUrl: './tags.component.html',
  styleUrls: ['./tags.component.scss'],
})
export class PostTagsComponent extends AbstractPostComponent { }
