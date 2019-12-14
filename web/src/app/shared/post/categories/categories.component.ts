import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';
import { PostComponent } from '../post.component';

@Component({
  selector: 'app-post-categories',
  changeDetection: ChangeDetectionStrategy.OnPush,
  templateUrl: './categories.component.html',
  styleUrls: ['./categories.component.scss'],
})
export class PostCategoriesComponent extends PostComponent implements OnInit {

  /**
   * List of categories
   */
  categories: Array<Category>;

  ngOnInit(): void {
    this.categories = this.post.categories;
  }

}
