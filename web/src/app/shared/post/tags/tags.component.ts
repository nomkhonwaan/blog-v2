import { ChangeDetectionStrategy, Component } from '@angular/core';
import { PostComponent } from '../post.component';

@Component({
  selector: 'app-post-tags',
  changeDetection: ChangeDetectionStrategy.OnPush,
  templateUrl: './tags.component.html',
  styleUrls: ['./tags.component.scss'],
})
export class PostTagsComponent extends PostComponent { }
