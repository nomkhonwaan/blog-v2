import { ChangeDetectionStrategy, Component, Input } from '@angular/core';
import { PostComponent } from '../post.component';

@Component({
  selector: 'app-post-author',
  changeDetection: ChangeDetectionStrategy.OnPush,
  templateUrl: './author.component.html',
  styleUrls: ['./author.component.scss'],
  styles: [
  ]
})
export class PostAuthorComponent extends PostComponent {

  /**
   * Use to indicate whether the display should show or not
   */
  @Input()
  displayName = true;

}
