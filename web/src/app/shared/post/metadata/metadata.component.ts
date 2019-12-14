import { ChangeDetectionStrategy, Component } from '@angular/core';
import { PostComponent } from '../post.component';

@Component({
  selector: 'app-post-metadata',
  changeDetection: ChangeDetectionStrategy.OnPush,
  templateUrl: './metadata.component.html',
  styleUrls: ['./metadata.component.scss'],
})
export class PostMetadataComponent extends PostComponent { }
